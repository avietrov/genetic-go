# Genetic algorithms in go

This project implements a [genetic algorithm](https://en.wikipedia.org/wiki/Genetic_algorithm) that approximates an image with a set of triangles.

## Examples
Using all the defaults to draw a gopher:

```
go run main.go --target="example/gopher.png"
```

<img src="example/gopher.png" width=200/> <img src="example/gopher.gif" width=200/>
A more complex image:
```
go run main.go --genom-size=1000 --mutation-count=30 --population-size=100 --x-over-count=4 --target="example/lisa.jpg"
```

<img src="example/lisa.jpg" width=200/> <img src="example/lisa.gif" width=200/>

## Configuration
You can configure the algorithm with the following parameters:
`genom-size` - how large is DNA of each individ (how many triangles per image). More triangles - more details, but also makes the algorithm slower.

`max-generations` - maximum number of generations, if the image is super simple the algorithm can converge and stop, but for any non-trivial image the algorithm stops by reaching maximum number of generations.

`mutation-count` - value from 1 to `genom-size` that defines how many genes are mutated per individ mutation.

`mutation-power` - value from 0.0+ to to 1.0, that defines stregth of mutations. 

`population-size` - number of individuals per generation (how many images per generation). 

`target` - target image (default "target.png")

`x-over-count` - number of crossovers done per generation. (default 2)

### How to choose a config
If your image has lots of details you should set `genom-size` to large values (check examples above for example values). With more mutations or stronger mutations (configured through `mutation-count` and `mutation-power` respectively) your learning process will stabilize faster, but you risk "jumping over" a better solution. I haven't found algorithmical advantages of different values here, since total amount of mutations and evaluations is `max-generations` * `x-over-count`, and you achive similar results by ballancing two values. Having said that, cross-overs are done in parralel, therefore I sent `x-over-count`  to number of cores in my CPU. 
