





84






# GF_ADMIN
/v1/admin/users/delete
/v1/admin/users/get_all
/v1/admin/users/get_all_invite_list
/v1/admin/users/add_to_invite_list
/v1/admin/users/remove_from_invite_list
/v1/admin/users/resend_confirm_email
/v1/admin/login_ui
/v1/admin/login
/v1/admin/dashboard
/v1/admin/healthz



# GF_ANALYTICS
/v1/a/ue
/v1/a/dashboard

# GF_CRAWL
/a/crawl/cluster/register__worker
/a/crawl/cluster/create__page_imgs
/a/crawl/cluster/create__page_img_ref
/a/crawl/cluster/link__get_unresolved
/a/crawl/cluster/link__mark_as_resolved
/a/crawl/image/recent
/a/crawl/image/add_to_flow
/a/crawl/search
/a/crawl/crawl_dashboard

# GF_DOMAINS
/a/domains/browser

# GF_HOME
/v1/home/viz/get
/v1/home/viz/update
/v1/home/view

# GF_IDENTITY
/v1/identity/eth/preflight
/v1/identity/eth/login
/v1/identity/eth/create
/v1/identity/userpass/login
/v1/identity/userpass/create
/v1/identity/policy/update
/v1/identity/email_confirm
/v1/identity/mfa_confirm
/v1/identity/update
/v1/identity/me
/v1/identity/register_invite_email

# GF_IMAGES - GF_GIF
/images/gif/get_info

# GF_IMAGES - GF_IMAGE_EDITOR
/images/editor/save

# GF_IMAGES - GF_IMAGES_FLOWS
/v1/images/flows/all
/v1/images/flows/add_img
/images/flows/add_img
/images/flows/imgs_exist
/images/flows/browser
/images/flows/browser_page

# GF_IMAGES - GF_IMAGES_JOBS
/images/jobs/start
/images/jobs/status

# GF_IMAGES
/v1/images/get
/images/d/
/v1/images/upload_init
/v1/images/upload_complete
/images/c
/images/v1/healthz

# GF_LANDING_PAGE
/landing/main

# GF_PUBLISHER
/posts/*
/posts/create
/posts/status
/posts/update
/posts/delete
/posts/browser
/posts/browser_page
/posts_elements/create

# GF_TAGGER
/v1/bookmarks/create
/v1/bookmarks/get
/v1/tags/notes/create
/v1/tags/notes/get
/v1/tags/create
/v1/tags/objects

# GF_WEB3 - GF_ADDRESS
/v1/web3/address/get_all
/v1/web3/address/add

# GF_WEB3 - GF_ETH_INDEXER
/gfethm/v1/block/index/job_updates
/gfethm/v1/block/index

# GF_WEB3 - GF_ETH_MONITOR_WORKER_INSPECTOR
/gfethm_worker_inspect/v1/account/info
/gfethm_worker_inspect/v1/tx/trace
/gfethm_worker_inspect/v1/blocks
/gfethm_worker_inspect/v1/health

# GF_WEB3 - GF_NFT
/v1/web3/nft/index_address
/v1/web3/nft/get_by_owner
/v1/web3/nft/get

# GF_WERB
/gfethm/v1/favorites/tx/add
/gfethm/v1/tx/trace/plot
/gfethm/v1/block
/gfethm/v1/miner
/gfethm/v1/peers
/gfethm/v1/health