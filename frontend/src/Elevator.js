class Elevator {
  constructor(x, y, width, height) {
    this.x = x;
    this.y = y;
    this.width = width;
    this.height = height;
    this.DoorOpen = false;
  }

  SetElevation(y) {
    this.y = y;
  }

  SetDoorOpen(value) {
    this.DoorOpen = value;
  }

  Draw() {
    if (this.DoorOpen) {
      fill(50, 150, 50);
    } else {
      fill(150, 50, 50);
    }

    line(this.x - this.width / 2 - 10, this.y, this.x + this.width / 2 + 10, this.y);
    rect(this.x, this.y, this.width, this.height);
  }
}

export default Elevator;
