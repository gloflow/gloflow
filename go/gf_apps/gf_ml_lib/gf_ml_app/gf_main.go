/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package main

import (
	"fmt"
	G "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
	"gopkg.in/cheggaaa/pb.v1"

	"time"
	"errors"
	"log"


	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------

type GFconvNet struct {
	g                  *G.ExprGraph
	w0, w1, w2, w3, w4 *G.Node // weights
	d0, d1, d2, d3     float64 // dropout probabilities
	out                *G.Node
}

//-------------------------------------------------

func main() {



	fmt.Println("YES")



	epochsInt := 100
	dataType := tensor.Float64
	batchSizeInt  := 32
	classesNumInt := 10
	examplesNumInt := 1000

	imageWidthInt    := 28
	imageHeightInt   := 28
	imagesInputShape := G.WithShape(batchSizeInt, 1, imageWidthInt, imageHeightInt)

	g := G.NewGraph()
	x := G.NewTensor(g, dataType, 4, imagesInputShape, G.WithName("x"))
	y := G.NewMatrix(g, dataType, G.WithShape(batchSizeInt, classesNumInt), G.WithName("y"))

	fmt.Println(g)


	model := newConvNet(g)
	fmt.Println(model)


	if err := fwd(x, model); err != nil {
		log.Fatalf("%+v", err)
	}


	losses := G.Must(G.Log(G.Must(G.HadamardProd(model.out, y))))
	cost := G.Must(G.Mean(losses))
	cost = G.Must(G.Neg(cost))
	// we wanna track costs
	var costVal G.Value
	G.Read(cost, &costVal)

	fmt.Println("loss:")
	spew.Dump(costVal)



	if _, err := G.Grad(cost, learnables(model)...); err != nil {
		log.Fatal(err)
	}


	prog, locMap, _ := G.Compile(g)
	vm := G.NewTapeMachine(g, G.WithPrecompiled(prog, locMap), G.BindDualValues(learnables(model)...))
	solver := G.NewRMSPropSolver(G.WithBatchSize(float64(batchSizeInt)))
	defer vm.Close()

	spew.Dump(solver)

	batchesNumInt := examplesNumInt / batchSizeInt
	log.Printf("Batches %d", batchesNumInt)

	bar := pb.New(batchesNumInt)
	bar.SetRefreshRate(time.Second)
	bar.SetMaxWidth(80)

	for i := 0; i < epochsInt; i++ {
		bar.Prefix(fmt.Sprintf("Epoch %d", i))
		bar.Set(0)
		bar.Start()


		for b := 0; b < batchesNumInt; b++ {

			start := b * batchSizeInt
			end := start + batchSizeInt
			if start >= examplesNumInt {
				break
			}
			if end > examplesNumInt {
				end = examplesNumInt
			}



			vm.Reset()
			bar.Increment()

		}


		log.Printf("Epoch %d | cost %v", i, costVal)
	}


}

func learnables(pModel *GFconvNet) G.Nodes {
	return G.Nodes{pModel.w0, pModel.w1, pModel.w2, pModel.w3, pModel.w4}
}

//-------------------------------------------------
func fwd(pX *G.Node,
	pModel *GFconvNet) error {

	var c0, c1, c2, fc *G.Node
	var a0, a1, a2, a3 *G.Node
	var p0, p1, p2 *G.Node
	var l0, l1, l2, l3 *G.Node

	var err error

	// LAYER 0
	// here we convolve with stride = (1, 1) and padding = (1, 1),
	// which is your bog standard convolution for convnet
	if c0, err = G.Conv2d(pX, pModel.w0, tensor.Shape{3, 3}, []int{1, 1}, []int{1, 1}, []int{1, 1}); err != nil {
		return errors.New("Layer 0 Convolution failed")
	}
	if a0, err = G.Rectify(c0); err != nil {
		return errors.New("Layer 0 activation failed")
	}
	if p0, err = G.MaxPool2D(a0, tensor.Shape{2, 2}, []int{0, 0}, []int{2, 2}); err != nil {
		return errors.New("Layer 0 Maxpooling failed")
	}
	log.Printf("p0 shape %v", p0.Shape())
	if l0, err = G.Dropout(p0, pModel.d0); err != nil {
		return errors.New("Unable to apply a dropout")
	}

	// Layer 1
	if c1, err = G.Conv2d(l0, pModel.w1, tensor.Shape{3, 3}, []int{1, 1}, []int{1, 1}, []int{1, 1}); err != nil {
		return errors.New("Layer 1 Convolution failed")
	}
	if a1, err = G.Rectify(c1); err != nil {
		return errors.New("Layer 1 activation failed")
	}
	if p1, err = G.MaxPool2D(a1, tensor.Shape{2, 2}, []int{0, 0}, []int{2, 2}); err != nil {
		return errors.New("Layer 1 Maxpooling failed")
	}
	if l1, err = G.Dropout(p1, pModel.d1); err != nil {
		return errors.New("Unable to apply a dropout to layer 1")
	}

	// Layer 2
	if c2, err = G.Conv2d(l1, pModel.w2, tensor.Shape{3, 3}, []int{1, 1}, []int{1, 1}, []int{1, 1}); err != nil {
		return errors.New("Layer 2 Convolution failed")
	}
	if a2, err = G.Rectify(c2); err != nil {
		return errors.New("Layer 2 activation failed")
	}
	if p2, err = G.MaxPool2D(a2, tensor.Shape{2, 2}, []int{0, 0}, []int{2, 2}); err != nil {
		return errors.New("Layer 2 Maxpooling failed")
	}

	var r2 *G.Node
	b, c, h, w := p2.Shape()[0], p2.Shape()[1], p2.Shape()[2], p2.Shape()[3]
	if r2, err = G.Reshape(p2, tensor.Shape{b, c * h * w}); err != nil {
		return errors.New("Unable to reshape layer 2")
	}
	log.Printf("r2 shape %v", r2.Shape())
	if l2, err = G.Dropout(r2, pModel.d2); err != nil {
		return errors.New("Unable to apply a dropout on layer 2")
	}

	// Layer 3
	if fc, err = G.Mul(l2, pModel.w3); err != nil {
		return errors.New("Unable to multiply l2 and w3")
	}
	if a3, err = G.Rectify(fc); err != nil {
		return errors.New("Unable to activate fc")
	}
	if l3, err = G.Dropout(a3, pModel.d3); err != nil {
		return errors.New("Unable to apply a dropout on layer 3")
	}

	// output decode
	var out *G.Node
	if out, err = G.Mul(l3, pModel.w4); err != nil {
		return errors.New("Unable to multiply l3 and w4")
	}
	pModel.out, err = G.SoftMax(out)
	
	return nil
}

//-------------------------------------------------
func newConvNet(pGraph *G.ExprGraph) *GFconvNet {

	dataType := tensor.Float64

	w0 := G.NewTensor(pGraph, dataType, 4, G.WithShape(32, 1, 3, 3), G.WithName("w0"), G.WithInit(G.GlorotN(1.0)))
	w1 := G.NewTensor(pGraph, dataType, 4, G.WithShape(64, 32, 3, 3), G.WithName("w1"), G.WithInit(G.GlorotN(1.0)))
	w2 := G.NewTensor(pGraph, dataType, 4, G.WithShape(128, 64, 3, 3), G.WithName("w2"), G.WithInit(G.GlorotN(1.0)))
	w3 := G.NewMatrix(pGraph, dataType, G.WithShape(128*3*3, 625), G.WithName("w3"), G.WithInit(G.GlorotN(1.0)))
	w4 := G.NewMatrix(pGraph, dataType, G.WithShape(625, 10), G.WithName("w4"), G.WithInit(G.GlorotN(1.0)))

	return &GFconvNet{
		g:  pGraph,
		w0: w0,
		w1: w1,
		w2: w2,
		w3: w3,
		w4: w4,

		d0: 0.2,
		d1: 0.2,
		d2: 0.2,
		d3: 0.55,
	}
}