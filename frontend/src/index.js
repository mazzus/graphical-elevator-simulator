import axios from "axios";
import queryString from "query-string";

import Lamp from "./Lamp";
import Button from "./Button";
import Elevator from "./Elevator";
import FloorIndicator from "./FloorIndicator";

let NFLOORS;

let backend;

let W;
let H;
let TOP;
let BOTTOM;

let elevatorCenter = 120;

let showDebug = false;

let elevator = null;

/**
 * @type {Button[]}
 */
let cabinButtons = [];
/**
 * @type {Button[]}
 */
let upButtons = [];
/**
 * @type {Button[]}
 */
let downButtons = [];
/**
 * @type {Button}
 */
let stopButton = null;
/**
 * @type {Button}
 */
let obstructionButton = null;

let floorIndicators = [];

let elev;

/**
 * @type {Lamp[]}
 */
let downLamps = [];

/**
 * @type {Lamp[]}
 */
let upLamps = [];

/**
 * @type {Lamp[]}
 */
let cabinLamps = [];

function FloorPosition(i) {
  return TOP + (BOTTOM - TOP) / (NFLOORS - 1) * (NFLOORS - i - 1);
}

window.setup = () => {
  const query = queryString.parse(location.search);

  backend = !query.backend ? "/api" : query.backend;
  W = !query.width ? 800 : query.width;
  H = !query.height ? 700 : query.height;
  NFLOORS = !query.floors ? 700 : query.floors;

  TOP = H / (2 * NFLOORS);
  BOTTOM = H - TOP;

  elev = new Elevator(elevatorCenter, 0, 50, H / (1.5 * NFLOORS));

  createCanvas(W, H);
  rectMode(CENTER);

  for (let i = NFLOORS; i > 0; i--) {
    cabinButtons.push(new Button(20, FloorPosition(NFLOORS - i), 30, 30, "none"));
  }

  for (let i = NFLOORS; i > 0; i--) {
    cabinLamps.push(new Lamp(55, FloorPosition(NFLOORS - i), 30, 30));
  }

  for (let i = NFLOORS; i > 1; i--) {
    upLamps.push(new Lamp(250, FloorPosition(NFLOORS - i), 30, 30));
  }

  for (let i = NFLOORS; i > 1; i--) {
    upButtons.push(new Button(290, FloorPosition(NFLOORS - i), 30, 30, "up"));
  }

  for (let i = NFLOORS - 1; i > 0; i--) {
    downButtons.push(new Button(350, FloorPosition(NFLOORS - i), 30, 30, "down"));
  }

  for (let i = NFLOORS - 1; i > 0; i--) {
    downLamps.push(new Lamp(390, FloorPosition(NFLOORS - i), 30, 30));
  }

  for (let i = 0; i < NFLOORS; i++) {
    floorIndicators.push(new FloorIndicator(200, FloorPosition(i) + 5, 20, i));
  }

  setInterval(fetchElevator, 12);
};

window.draw = () => {
  background(190);
  if (!elevator) {
    console.warn("no elevator");
    return;
  }

  for (let i = 0; i < NFLOORS; i++) {
    let linePosUpper = FloorPosition(i + elevator.margin);
    let linePosLower = FloorPosition(i - elevator.margin);
    fill(100, i === elevator.currentFloor ? 255 : 0, 50);
    rect(elevatorCenter, (linePosUpper + linePosLower) / 2, 70, linePosUpper - linePosLower);
  }

  // Draw Cabin box
  fill(70, 70, 70);
  rect(40, H / 2, 80, H);
  fill(255);
  textSize(20);
  text("CABIN", 10, 30);

  // Draw Up Box
  fill(70, 70, 70);
  rect(270, H / 2, 80, H);
  fill(255);
  textSize(20);
  text("UP", 250, 30);

  // Draw Down box
  fill(70, 70, 70);
  rect(370, H / 2, 80, H);
  fill(255);
  textSize(20);
  text("DOWN", 340, 30);

  if (showDebug) {
    fill(0);
    textSize(10);
    textStyle(NORMAL);
    text(JSON.stringify(elevator, null, 2), 430, 10);
  }

  cabinButtons.forEach((b, i) => {
    b.SetRegistered(elevator.cabinButtons[i]);
    b.Draw();
  });

  upButtons.forEach((b, i) => {
    b.SetRegistered(elevator.upButtons[i]);
    b.Draw();
  });

  downButtons.forEach((b, i) => {
    b.SetRegistered(elevator.downButtons[i + 1]);
    b.Draw();
  });

  cabinLamps.forEach((l, i) => {
    l.SetLight(elevator.cabinLamps[i]);
    l.Draw();
  });

  downLamps.forEach((l, i) => {
    l.SetLight(elevator.downLamps[i + 1]);
    l.Draw();
  });

  upLamps.forEach((l, i) => {
    l.SetLight(elevator.upLamps[i]);
    l.Draw();
  });

  floorIndicators.forEach((l, i) => {
    l.SetActive(i == elevator.indicatorLamp);
    l.Draw();
  });

  elev.SetElevation(FloorPosition(elevator.position));
  elev.SetDoorOpen(elevator.doorLamp);
  elev.Draw();
};

window.mousePressed = () => {
  cabinButtons.forEach((b, i) => {
    if (b.Press(mouseX, mouseY)) {
      SendButton("cabin", i, true);
    }
  });

  upButtons.forEach((b, i) => {
    if (b.Press(mouseX, mouseY)) {
      SendButton("up", i, true);
    }
  });

  downButtons.forEach((b, i) => {
    if (b.Press(mouseX, mouseY)) {
      SendButton("down", i + 1, true);
    }
  });
};

function SendButton(type, floor, value) {
  axios.post(backend + "/button", { Type: type, Floor: floor, Value: value });
}

window.mouseReleased = () => {
  cabinButtons.forEach((b, i) => {
    if (b.Release(mouseX, mouseY)) {
      SendButton("cabin", i, false);
    }
  });

  upButtons.forEach((b, i) => {
    if (b.Release(mouseX, mouseY)) {
      SendButton("up", i, false);
    }
  });

  downButtons.forEach((b, i) => {
    if (b.Release(mouseX, mouseY)) {
      SendButton("down", i + 1, false);
    }
  });
};

window.keyPressed = () => {
  if (keyCode === 68) {
    // "D"
    showDebug = !showDebug;
  }
};

function fetchElevator() {
  axios
    .get(backend + "/total")
    .then(response => {
      elevator = response.data;
    })
    .catch(err => {
      console.error({ err });
    });
}
