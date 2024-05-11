import * as Plot from "npm:@observablehq/plot";

export function piano_roll(notes, {width, height} = {}) {
  const list = [
    [0,   10],[0.5,  10],[null,10],
    [0.5, 20],[0.75, 20],[null,20],
    [1.0, 30],[1.75, 30],[null,30],
    [1.25,40],[1.75, 40],[null,40],
    [1.75,50],[2.5,  50],[null,50],
    [3.5, 60],[3.75, 60],[null,60],
  ]

  return Plot.plot({
    width,
    height,
    marginTop: 30,
    x: {nice: true, label: null, tickFormat: ""},
    y: {axis: null},
    marks: [
      Plot.line(list, {
              stroke: "black",
              strokeWidth: 4
            }),
      Plot.rect(notes, {
              x1: "start",
              y1: "y1",
              x2: "end",
              y2: "y2",
              stroke: "red",
              strokeWidth: 1
            }),
      // Plot.ruleX(notes, {x: "start", y: "y", markerEnd: "dot", strokeWidth: 1.5}),
      Plot.ruleY([0]),
      Plot.text(notes, {x: "start", y: "y", text: "note", lineAnchor: "bottom", dy: -10, lineWidth: 10, fontSize: 12})
    ]
  });

}
