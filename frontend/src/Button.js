class Button {
  constructor(x, y, width, height, type) {
    this.x = x;
    this.y = y;
    this.width = width;
    this.height = height;
    this.registered = false;
    this.type = type;
  }

  inBounds(x, y) {
    return (
      x > this.x - this.width / 2 && x < this.x + this.width / 2 && y > this.y - this.width / 2 && y < this.y + this.width / 2
    );
  }

  Press(x, y) {
    return this.inBounds(x, y);
  }

  Release(x, y) {
    return this.inBounds(x, y);
  }

  Draw() {
    if (this.registered) {
      fill(30, 30, 30);
    } else {
      fill(100, 100, 100);
    }

    if (this.type === "up") {
      triangle(
        this.x - this.width / 2,
        this.y + this.height / 2,
        this.x + this.width / 2,
        this.y + this.height / 2,
        this.x,
        this.y - this.height / 2
      );
    } else if (this.type === "down") {
      triangle(
        this.x - this.width / 2,
        this.y - this.height / 2,
        this.x,
        this.y + this.height / 2,
        this.x + this.width / 2,
        this.y - this.height / 2
      );
    } else if (!this.type || this.type == "none") {
      rect(this.x, this.y, this.height, this.width);
    }
  }

  SetRegistered(value) {
    this.registered = value;
  }
}

export default Button;
