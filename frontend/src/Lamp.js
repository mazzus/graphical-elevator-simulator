class Lamp {
  constructor(x, y, width, height) {
    this.x = x;
    this.y = y;
    this.width = width;
    this.height = height;
    this.on = false;
  }

  /**
   *
   * @param {Boolean} value
   */
  SetLight(value) {
    this.on = value;
  }

  Draw() {
    if (this.on) {
      fill(200, 0, 0);
    } else {
      fill(100, 100, 100);
    }
    ellipse(this.x, this.y, this.width, this.height);
  }
}

export default Lamp;
