import os
import json
from typing import Any, Dict
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


def create_flow(d: WebDriver, w: WebDriverWait, info: Dict[Any, Any]):

  # Go to recipe list screen
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img")).click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-recipe img").click()

  # Create recipe
  d.find_element(By.ID, "recipe-index-create-btn").click()

  # Name
  d.find_element(
      By.CSS_SELECTOR, ".recipe-create-name input").send_keys(info["name"])
  d.find_element(By.CSS_SELECTOR, ".recipe-create-submit").click()

  # Meal
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".recipe-single-meal select")).send_keys(info["meal"])

  # Description
  d.find_element(
      By.CSS_SELECTOR, ".recipe-single-desc textarea").send_keys(info["desc"])

  # Make ingredients editable
  d.find_element(
      By.CSS_SELECTOR, ".recipe-single-ingredient-edit img").click()
  # Add ingredient
  add_ingredient = d.find_element(
      By.CSS_SELECTOR, ".recipe-single-ingredient-add img")

  for i, ingre in enumerate(info["ingre"]):
    add_ingredient.click()
    try: 
      w.until(lambda x: len(x.find_elements(
          By.CLASS_NAME, "recipe-ingredient-root")) > i)
    except Exception as e:
      raise e
    ingre_ele = d.find_elements(By.CLASS_NAME, "recipe-ingredient-root")[i]
    # Add ingredients
    ingre_ele.find_element(
        By.CSS_SELECTOR,
        ".recipe-ingredient-quantity input").send_keys(Keys.BACKSPACE + ingre["quantity"])
    ingre_ele.find_element(
        By.CSS_SELECTOR, ".recipe-ingredient-unit").send_keys(ingre["unit"])
    item = ", ".join(
      [ingre["item"], ingre["modifier"]]) if "modifier" in ingre else ingre["item"]
    ingre_ele.find_element(
        By.CSS_SELECTOR, ".recipe-ingredient-item input").send_keys(item)
    ingre_ele.find_element(
        By.CSS_SELECTOR, ".recipe-ingredient-item input").send_keys(Keys.TAB)
    if ingre.get("optional", False):
      ingre_ele.find_element(
          By.CSS_SELECTOR, ".recipe-ingredient-optional input").click()

  d.find_element(By.CSS_SELECTOR, "ol.recipe-instruction-list").click()
  instr = d.find_element(By.CSS_SELECTOR, ".recipe-instruction-list textarea")
  instr.send_keys(info["instr"])
  instr.send_keys(Keys.TAB)  # Unfocus for verify

  verify(d, info)
  # Save
  w.until(lambda x: x.find_element(
      By.XPATH, "//div[contains(@class,'recipe-single-name')]/*/input[@value!='']"))
  verify(d, info)

  # Go to recipe list screen
  d.refresh()
  w.until(lambda x: x.find_element(
      By.XPATH, "//div[contains(@class,'recipe-single-name')]/*/input[@value!='']"))
  verify(d, info)
