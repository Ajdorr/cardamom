from typing import Any
from tests import login, clear
from . import verify
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys

info = {
    "name": "Chocolate Chip Cookies",
    "desc": "Buttery and irresistible chocolate chip cookies",
    "meal": "Dessert",
    "ingre": [
        {
            "quantity": "1/2",
            "unit": "cup",
            "item": "flour"
        },
        {
            "quantity": "1/2",
            "unit": "cup",
            "item": "butter"
        },
        {
            "quantity": "1/4",
            "unit": "cup",
            "item": "chocolate chips"
        },
        {
            "quantity": "1",
            "unit": "none",
            "item": "egg"
        },
        {
            "quantity": "1/4",
            "unit": "cup",
            "item": "milk"
        }
    ],
    "instr": "Combine pancake mix and water, whisk until lumps are gone\ngrease pan and pour approx 1 Tbsp per pancake\nflip after approximately 5 minutes, then for an additional 3 minutes or until golden brown"

}


def test_complex():

  d, w = login()

  # Go to recipe list screen
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img")).click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-recipe img").click()

  # Create recipe
  d.find_element(By.ID, "recipe-index-create-btn").click()

  # Name
  name_ele = d.find_element(By.CSS_SELECTOR, ".recipe-single-name input")
  name_ele.send_keys(info["name"][:-4])
  clear(d, name_ele)
  name_ele.send_keys(info["name"])

  # Meal
  meal_ele = d.find_element(By.CSS_SELECTOR, ".recipe-single-meal select")
  meal_ele.send_keys(info["meal"])

  # Description
  desc_ele = d.find_element(By.CSS_SELECTOR, ".recipe-single-desc textarea")
  desc_ele.send_keys(info["desc"][3:-4])
  clear(d, desc_ele)
  desc_ele.send_keys(info["desc"])

  # Add ingredient
  add_ingre = d.find_element(
      By.CSS_SELECTOR, ".recipe-single-ingredient-add img")
  for _ in info["ingre"]:
    add_ingre.click()

  for i, ingre in enumerate(d.find_elements(By.CLASS_NAME, "recipe-ingredient-root")):
    quantity = ingre.find_element(
        By.CSS_SELECTOR, ".recipe-ingredient-quantity input")
    unit = ingre.find_element(By.CSS_SELECTOR, ".recipe-ingredient-unit")
    item = ingre.find_element(By.CSS_SELECTOR, ".recipe-ingredient-item input")

    quantity.send_keys("zz")
    unit.send_keys("none")
    item.send_keys(info["ingre"][i]["item"][1])

    clear(d, quantity)
    quantity.send_keys(info["ingre"][i]["quantity"])
    unit.send_keys(info["ingre"][i]["unit"])
    clear(d, item)
    item.send_keys(info["ingre"][i]["item"])

  d.find_element(By.CSS_SELECTOR, "ol.recipe-instruction-list").click()
  instr = d.find_element(By.CSS_SELECTOR, ".recipe-instruction-list textarea")
  instr.send_keys("test data")
  clear(d, instr)
  instr.send_keys(info["instr"])
  instr.send_keys(Keys.TAB)

  verify(d, info)
  # Save
  d.find_element(By.CSS_SELECTOR, ".recipe-single-save img").click()
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img"))
  verify(d, info)
