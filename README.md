

![alt text](http://gf--img.s3-website-us-east-1.amazonaws.com/gf_logo_0.3.png "GloFlow logo")

Media publishing/management/manipulation system.

(note: this is alpha software, not feature complete, and still missing more code that needs to be migrated to this repo)

Its purpose is to help manage media (for now mainly images) libraries, to **edit** media, **publish** it, **share** it with other people, to **analyse** it. 
Its goal is to provide free and private exchange of media between people and groups of people. It currently runs as a centralized service, but my aim is for it to also be fully decentralized and run in a P2P network. 

Images/videos are at the core of our **culture**, at the core of how we perceive the world and how we **remember ideas** and **moments**. **We communicate most effectivelly when we exchange images**. There should be a technology that is focused on that, that is modern, accessible to technical people for modification/integration and automation, and most importantly **free** and **independent** of any single individuals or groups (and their possible control).

This project is still very much work in progress. I doubt that it will ever be done, or truly "ready". I worked on it over the years in an largely unplanned manner, working on parts that were interesting to me, when I had free time and wanted to just hack without deadlines or immediate purpose. There were various pauses in development due to work and general life, but I always came back to it. 
It has been rewritten 3 times since I wrote the first official code in **2012** (that went online). 

Originally it was coded in Python on the backend, and JS on the frontend. Later the entire thing was rewritten in Dart on both the frontend and backend. Finally the backend was rewritten in **Go**, and frontend in TypeScript. 



A single style is maintained across languages used in the implementation (**Go**,**Python**,**Typescript**) - even though the languages are different enough from each other. The focus is on basic functional language principles (of pure functions, high-level functions, closures). Functions should receive all the state that they operate on via their arguments (other then functions that work with external state - DB or external queries). Object orientation (objects holding state and methods operating on that state internally) is avoided as much as possible (even though it is the default idiomatic style of Go and Python). State/variable mutation still exists in various places, but the aim is to keep it to a minimum (constant runtime values would be a welcome feature in Go and Python). 




Originally created by Ivan Trajkovic
