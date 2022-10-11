import os
import json
from typing import Any
from tests import login
from . import verify
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.remote.webdriver import WebDriver
from selenium.webdriver.support.wait import WebDriverWait


def test_creates():
  fp = os.path.join(
      os.path.dirname(__file__),
      "../resources/recipe/test_data.json")

  d, w = login()
  cases = json.load(open(fp))
  for c in cases:
    create_flow(d, w, c)


def create_flow(d: WebDriver, w: WebDriverWait, info: dict[Any, Any]):

  # Go to recipe list screen
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img")).click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-recipe img").click()

  # Create recipe
  d.find_element(By.ID, "recipe-index-create-btn").click()

  # Name
  d.find_element(
      By.CSS_SELECTOR, ".recipe-single-name input").send_keys(info["name"])

  # Meal
  d.find_element(
      By.CSS_SELECTOR, ".recipe-single-meal select").send_keys(info["meal"])

  # Description
  d.find_element(
      By.CSS_SELECTOR, ".recipe-single-desc textarea").send_keys(info["desc"])

  # Add ingredient
  add_ingredient = d.find_element(
      By.CSS_SELECTOR, ".recipe-single-ingredient-add img")
  for _ in info["ingre"]:
    add_ingredient.click()

  ingres = d.find_elements(By.CLASS_NAME, "recipe-ingredient-root")
  for i, ingre in enumerate(info["ingre"]):
    # Add ingredients
    ingres[i].find_element(
        By.CSS_SELECTOR,
        ".recipe-ingredient-quantity input").send_keys(Keys.BACKSPACE + ingre["quantity"])
    ingres[i].find_element(
        By.CSS_SELECTOR, ".recipe-ingredient-unit").send_keys(ingre["unit"])
    ingres[i].find_element(
        By.CSS_SELECTOR, ".recipe-ingredient-item input").send_keys(ingre["item"])

  add_instr = d.find_element(
      By.CSS_SELECTOR, ".recipe-instruction-input input")
  for instr in info["instr"]:
    add_instr.send_keys(instr)
    add_instr.send_keys(Keys.ENTER)

  verify(d, info)
  # Save
  d.find_element(By.CSS_SELECTOR, ".recipe-single-save img").click()
  w.until(lambda x: x.find_element(
      By.XPATH, "//div[contains(@class,'recipe-single-name')]/*/input[@value!='']"))
  verify(d, info)

  # Go to recipe list screen
  d.refresh()
  w.until(lambda x: x.find_element(
      By.XPATH, "//div[contains(@class,'recipe-single-name')]/*/input[@value!='']"))
  verify(d, info)
