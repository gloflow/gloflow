# GCN research (initial)

Initial research notes on Graph Convolutional Networks (GCNs): classification, scale, latency, low-memory viability and graph generation.
Each section is one question that was asked, followed by the summarized answer.
Figures are order-of-magnitude estimates from typical setups, not benchmarks of our own data. Profile on real data before relying on them.

## Contents
1. [Review of a pasted GCN overview](#1-review-of-a-pasted-gcn-overview)
2. [Can GCNs do classification?](#2-can-gcns-do-classification)
3. [How large can graphs be?](#3-how-large-can-graphs-be-loaded-and-classified)
4. [Inference latency for classification](#4-inference-latency-for-classifying-a-graph)
5. [Viability on old, low-memory GPUs](#5-viability-on-old-low-memory-gpus)
6. [Can GCNs do graph generation?](#6-can-gcns-do-graph-generation)
7. [Latency of generative models for classification](#7-latency-of-generative-models-for-classification)
8. [Open questions / next steps](#8-open-questions--next-steps)

---

## 1. Review of a pasted GCN overview

**Q: "how familiar are you with GCN's?"** (a pasted overview text on encoding graphs for CNN-style training was included)

**Core formula (Kipf & Welling):**

```
H^(l+1) = σ( D̃^(-1/2) Ã D̃^(-1/2) H^(l) W^(l) ),   Ã = A + I,   D̃_ii = Σ_j Ã_ij
```

- `X ∈ R^(N×F)`: node feature matrix. `A ∈ R^(N×N)`: adjacency matrix. Optional edge features `E ∈ R^(N×N×D)`.
- Node tasks: per-node softmax with cross-entropy over labeled nodes.
- Graph tasks: permutation-invariant readout (mean/sum/max/attention/DiffPool), then an MLP.
- Embeddings: unsupervised via VGAE (link prediction) or contrastive via DGI.

**Accurate in the pasted text:** the propagation rule, the node-level and graph-level pipelines, and the VGAE / DGI / DiffPool / PATCHY-SAN descriptions.

**Imprecise or misleading in the pasted text:**
1. The GCN layer is a first-order Chebyshev approximation of a spectral filter, plus the renormalization trick (self-loops before normalizing) for numerical stability.
2. Labelling the formula "spatial / message passing" and "first-order spectral" at once is contradictory. It is spectral in origin and *interpretable* as message passing.
3. Symmetric normalization is not a plain mean. Edges are weighted `1/sqrt(d_i d_j)`. A true mean aggregator uses `D^-1 A`.
4. Symmetric normalization assumes an undirected graph. Directed graphs need in/out-degree variants (DiGCN, MagNet).
5. The link-prediction loss is loosely written. The standard form is negative sampling.
6. The attention pooling formula needs the softmax normalized across nodes, with `MLP(Z_i)` giving a scalar (or per-channel) score.
7. "GCN/GraphSAGE/GIN significantly outperform grid adaptations" is unsupported as stated. Also, GCN is the weakest of the three in expressiveness (GIN matches 1-WL, plain GCN does not).

**Missing from the pasted text, but important in practice:**
- **Over-smoothing:** stacking more than 2-4 layers makes node representations converge. Fixes: residuals, JKNet, PairNorm, DropEdge.
- **Over-squashing:** long-range information is compressed through bottleneck nodes.
- **Scalability:** full-batch GCN needs the whole graph in memory. Use neighbor sampling (GraphSAGE), Cluster-GCN or GraphSAINT.
- **Transductive vs inductive:** the original paper is transductive. Weights are shared, so it can be applied inductively, but the normalization uses global graph structure.
- **Edge features:** plain GCN ignores them. Use R-GCN, MPNN, GINE or GAT with edge attributes.
- **Libraries:** PyTorch Geometric (`GCNConv`) and DGL.

---

## 2. Can GCNs do classification?

**Q: "so these networks allow for classification?"**

Yes, at three levels:

1. **Node classification.** Predict a label per node (paper topic, bot account). Usually semi-supervised: train on few labeled nodes, structure propagates information. Output `softmax(Z W)`, cross-entropy over labeled nodes.
2. **Graph classification.** One label per graph (molecule toxicity, document category). GCN layers, then pooling (mean/sum/max/attention), then MLP + softmax.
3. **Edge classification / link prediction.** Combine endpoint embeddings (`Z_i · Z_j` or `MLP([Z_i, Z_j])`) and classify.

**Why not a plain MLP on node features?** Each layer mixes a node's features with its neighbors', so the classifier sees the node's attributes plus local graph context. This helps most under **homophily** (connected nodes share labels). Under **heterophily**, plain GCN can do worse than an MLP. H2GCN and GPR-GNN handle that case better.

**Minimal PyG example (node classification):**

```python
import torch, torch.nn.functional as F
from torch_geometric.nn import GCNConv

class GCN(torch.nn.Module):
    def __init__(self, in_dim, hid, n_classes):
        super().__init__()
        self.c1 = GCNConv(in_dim, hid)
        self.c2 = GCNConv(hid, n_classes)

    def forward(self, x, edge_index):
        x = F.relu(self.c1(x, edge_index))
        x = F.dropout(x, p=0.5, training=self.training)
        return self.c2(x, edge_index)  # logits

out = model(data.x, data.edge_index)
loss = F.cross_entropy(out[data.train_mask], data.y[data.train_mask])
```

For graph classification, add `global_mean_pool(x, batch)` after the conv layers and a linear head on top.

---

## 3. How large can graphs be (loaded and classified)?

**Q: "how large are graphs that can be loaded and classified?"**

The limit is almost always memory, not the model.

**Full-batch node classification (whole graph on one device).** Memory is roughly:
- Edges: `2 x E x 8 bytes` for an int64 `edge_index` (100M edges is about 1.6 GB).
- Features: `N x F x 4 bytes`.
- Activations: about `L x N x d x 4 bytes`, times about 2-3 during training for backprop.

| GPU memory | Approx. size that fits (hidden 256, 2 layers) |
|---|---|
| 8-16 GB | ~100K-1M nodes, ~10-50M edges |
| 40-80 GB | ~2-5M nodes, ~100M+ edges |

**Reference points:**
- Cora: 2.7K nodes (trivial).
- Reddit: 233K nodes, 114M edges (full-batch feasible on a single large GPU).
- ogbn-products: 2.4M nodes, 62M edges (full-batch possible on 80 GB, usually sampled).
- ogbn-papers100M: 111M nodes, 1.6B edges (needs sampling plus hundreds of GB CPU RAM, or distributed).
- Pinterest PinSage: ~3B nodes, 18B edges (distributed).

**Larger graphs, without loading the whole graph on the GPU:**
- **Neighbor sampling** (GraphSAGE, PyG `NeighborLoader`): graph stays in CPU RAM, so RAM is the limit instead of VRAM.
- **Subgraph batching** (Cluster-GCN, GraphSAINT).
- **Precomputed propagation** (SGC, SIGN): compute `Â^k X` offline, then train a plain MLP.
- **Distributed training** (DGL-Dist, multi-GPU PyG, GraphLearn).

**Graph classification:** individual graphs are usually small (molecules, 20-100 nodes). Scale comes from the number of graphs, packed into a block-diagonal batch. Single huge graphs (100K+ nodes) need the techniques above.

**Dense adjacency is the real limit:** a dense `N x N` matrix tops out around 20-50K nodes (50K nodes in float32 is 10 GB). Use sparse formats (COO/CSR) for millions of nodes.

**Depth costs memory too:** with sampling, k-hop neighborhoods grow exponentially (neighbor explosion). Fan-outs like `[15, 10]` keep it manageable.

---

## 4. Inference latency for classifying a graph

**Q: "once trained what is the latency for classifying a graph?"**

**Cost model:** one GCN layer is about `O(E*d + N*d^2)` (sparse aggregation plus dense weight multiply). Inference is one forward pass, no backward pass.

**Graph classification (small graphs, e.g. molecules):**

| Setup | Typical latency |
|---|---|
| One graph (~30 nodes), CPU, PyG | ~1-5 ms |
| One graph, GPU | ~1-3 ms (kernel launch overhead dominates) |
| Batched, GPU | ~10K-100K+ graphs/s (microseconds per graph) |

For tiny graphs, framework overhead (Python, kernel launches, data conversion) exceeds the math. ONNX/TorchScript export or `torch.compile` reduces it.

**Node classification on a large graph:**
- Full-graph GPU pass: ~10-100 ms for 100K-1M nodes. Reddit scale (114M edges): ~100 ms to a few seconds.
- Single-node prediction with k-hop sampling: ~5-50 ms (dominated by neighborhood fetch and subgraph build).
- Precomputed embeddings / propagation (SIGN, SGC): sub-millisecond (MLP lookup), only for static graphs.

**Where the time goes in real systems:**
1. Graph construction and feature loading (often more than the forward pass).
2. CPU to GPU transfer for small requests.
3. Neighbor explosion on high-degree graphs (cap with fan-out limits).
4. Python overhead when serving single small graphs.

**Ways to cut latency:**
- Fewer layers (2 is usually enough), smaller hidden size.
- Dynamic batching of requests.
- Cache embeddings for nodes whose neighborhoods did not change.
- INT8 quantization, TensorRT / ONNX Runtime.
- Distill the GNN into an MLP (GLNN): large speedup, modest accuracy loss on many datasets.

Best practice: profile a forward pass on real data (`torch.utils.benchmark` or `torch.profiler`, ~10 lines).

---

## 5. Viability on old, low-memory GPUs

**Q: "for achieving super low memory consumption on old T40 type GPU's, are GCN's viable or obsoleted? worth studying?"**

"T40" is ambiguous: assumed Tesla T4, P40 or K40.

| GPU | Arch | Notes |
|---|---|---|
| **T4** (16 GB) | Turing | Tensor cores, good FP16/INT8. Fully supported by current PyTorch/PyG. |
| **P40** (24 GB) | Pascal | FP16 very slow (~1/64 rate), so use FP32 or INT8. Support is fading. |
| **K40** (12 GB) | Kepler | Compute capability 3.5. Recent PyTorch builds dropped it. Avoid if possible. |

**Verdict: viable, and among the lightest graph models. Not state of the art for accuracy, but far from obsolete as a baseline.**

**Why GCN suits low memory:**
- No attention, so no per-edge attention tensor (GAT stores `E x heads` values).
- Normalization `D̃^-1/2 Ã D̃^-1/2` can be precomputed once.
- Tiny parameter count (`F x d` per layer, often under 1 MB total).
- Memory is dominated by `N x d` activations, controllable through hidden size and depth.

**Low-memory techniques:**
- CSR/sparse format with **int32** indices instead of int64 (halves edge storage).
- Small hidden size (16-64) and 2 layers (the original Kipf model used 16).
- **SGC / SIGN**: precompute `Â^k X`, then a plain linear model or MLP.
- Sampling or Cluster-GCN to bound the working set.
- FP16 on T4, INT8 for inference.
- Gradient checkpointing for deeper models.
- Distill into an MLP (GLNN) so inference needs no graph.
- Avoid dense adjacency.

**Is it obsolete?**
- Accuracy: newer models beat it on many benchmarks, but several papers show a well-tuned GCN matches or beats fancier models on common node-classification sets.
- Efficiency: GCN still wins under tight memory budgets.
- Limits: weak on heterophily, long-range dependencies and expressiveness (below 1-WL), plus over-smoothing limits depth.

**Worth studying: yes.**
1. Almost every GNN is a message-passing variation, and GCN is the cleanest instance.
2. It is the standard baseline for judging other models.
3. Next to learn: **GraphSAGE, GIN, SGC/SIGN**, which cover the scalability and expressiveness gaps.

---

## 6. Can GCNs do graph generation?

**Q: "can they do graph generation?"**

Not by themselves. A GCN is an **encoder**, with no mechanism for producing new graphs. It is a common building block in generative models.

| Model | How it works | Limits |
|---|---|---|
| **GAE / VGAE** | GCN encoder gives `Z`. Decoder `σ(Z Z^T)` guesses the adjacency. VGAE adds a latent distribution so `Z` can be sampled. | Edges only, for a fixed node set. Poor samples, weak decoder. |
| **GraphVAE** | GCN encoder to one latent vector. MLP decoder outputs whole `A` and `X`. | Fixed max size, `O(N^2)` output, node-matching problem in the loss. |
| **MolGAN** | GAN with an R-GCN discriminator. Generator outputs dense adjacency/feature tensors. | Mode collapse, small molecules only. |
| **GraphRNN / GRAN** | Autoregressive: add nodes and edges one at a time (GRAN uses GNN layers per step). | Slow sampling, ordering matters. |
| **JT-VAE** | Builds molecules from substructures via a junction tree. GNNs encode the pieces. | Chemistry-specific. |
| **Graph diffusion** (DiGress, GDSS, EDP-GNN) | Add noise to the graph, train a GNN / graph transformer to denoise step by step. | Best quality currently, but costly (many denoising steps). |
| **Normalizing flows** (GraphNVP, MoFlow) | Invertible GNN-based layers map graph to/from a Gaussian. | Fixed size, complex. |

**Why generation is harder than classification:**
- **Discrete structure:** edges are 0/1, so gradients do not flow through sampling. Solutions: continuous relaxation, Gumbel-softmax, autoregressive or diffusion approaches.
- **Permutation invariance:** the same graph has `N!` adjacency matrices.
- **Variable size:** most decoders assume a fixed max `N`.
- **Validity:** e.g. chemical valence rules, enforced post-hoc or during sampling.

**Practical guidance:**
- Small graphs such as molecules: use diffusion (DiGress) or an autoregressive model, which give much better samples than VGAE.
- Predicting missing edges in an existing graph: VGAE (link prediction, not true generation) works well.
- GCN layers can serve as backbone, but graph diffusion models often swap them for graph transformers or GINE-style layers for more expressiveness.
- On a small GPU, GraphVAE-style models are the cheapest to train. Diffusion needs much more compute.

---

## 7. Latency of generative models for classification

**Q: "what is the latency of it? for classification"** (ambiguous: interpreted as generative models used for classification)

Usually only the encoder is kept at inference, so the generative part adds no cost.

| Model | Classification latency | Why |
|---|---|---|
| **VGAE / GAE encoder + classifier head** | Same as GCN: ~1-5 ms per small graph, 10-100 ms for large graphs on GPU | Decoder dropped, one forward pass through 2 GCN layers. |
| **GraphVAE encoder** | Same as GCN | Same reasoning. |
| **GraphRNN / GRAN** | ~10 ms to seconds per graph | Autoregressive. Scoring the likelihood goes node by node, so cost grows with `N`. |
| **Diffusion (DiGress etc.) as classifier** | ~100 ms to tens of seconds per graph | Each denoising step is a full pass, 50-1000 steps. Rarely used for classification. |
| **Flow models** | ~5-50 ms | One invertible pass, heavier than a GCN. |

**Generation latency (generative direction), for comparison:**
- VGAE / GraphVAE: ~ms (single decode).
- GRAN / GraphRNN: tens of ms to seconds, depending on size.
- Diffusion: seconds per graph (a 30-node molecule, 500 steps, T4: about 5-30 s).

**Takeaway:**
- If the only goal is classification, use a plain GCN, GIN or GraphSAGE (fastest).
- Generative pretraining (VGAE, DGI) can help when labels are few, and costs only GCN-level latency at inference (encoder only).
- A heavy generative model as the classifier itself is rarely worth the latency.

---

## 8. Open questions / next steps

To pick a concrete approach, the following still need to be defined:
- Task: node, graph or edge classification, or generation? What are the labels?
- Graph size: nodes and edges, static or dynamic graph, single large graph or many small ones?
- Target hardware: exact GPU (T4 / P40 / K40) and memory budget.
- Latency target: per-request (online) or batch (offline)?
- Homophily: do connected nodes share labels? (decides plain GCN vs heterophily-aware models)

Suggested learning order: GCN (Kipf & Welling), then GraphSAGE, GIN, SGC/SIGN, GLNN (distillation), then VGAE / DGI, then graph diffusion (DiGress) if generation is needed.

Key references: Kipf & Welling (GCN, GAE/VGAE), Hamilton et al. (GraphSAGE), Xu et al. (GIN), Wu et al. (SGC), Frasca et al. (SIGN), Chiang et al. (Cluster-GCN), Zeng et al. (GraphSAINT), Zhang et al. (GLNN), Velickovic et al. (DGI), Ying et al. (DiffPool), Vignac et al. (DiGress), Simonovsky & Komodakis (GraphVAE), De Cao & Kipf (MolGAN).
