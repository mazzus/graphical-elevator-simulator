class FloorIndicator {
  constructor(x, y, size, floor) {
    this.x = x;
    this.y = y;
    this.size = size;
    this.floor = floor;

    this.active = false;
  }

  Draw() {
    textSize(this.size);
    textStyle(BOLD);

    if (this.active) {
      fill(50, 200, 50);
    } else {
      fill(20, 20, 20);
    }

    text(`${this.floor}`, this.x, this.y);
  }
  SetActive(value) {
    this.active = value;
  }
}

export default FloorIndicator;
