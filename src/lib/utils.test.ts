import { add } from "./utils";
import { expect } from "@jest/globals";

describe("add", () => {
  it("should add two numbers", () => {
    expect(add(1, 2)).toBe(3);
  });
});
