from . import verify
from tests import login, clear
from tests.recipe.test_create import create_flow
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys

before_img = {
    "name": "Sandwich",
    "desc": "sandwich",
    "meal": "Breakfast",
    "ingre": [
        {
            "quantity": "2",
            "unit": "none",
            "item": "bread"
        },
        {
            "quantity": "50",
            "unit": "mL",
            "item": "ketchup"
        },
        {
            "quantity": "150",
            "unit": "g",
            "item": "ham"
        },
        {
            "quantity": "150",
            "unit": "g",
            "item": "cheese"
        },
    ],
    "instr": "spread ketchup on bread\nadd cheese and ham"

}
after_img = {
    "name": "Gourmet ham and cheese sandwich",
    "desc": "A sandwich fit for royalty",
    "meal": "Lunch",
    "ingre": [
        {
            "quantity": "2",
            "unit": "none",
            "item": "bread"
        },
        {
            "quantity": "50",
            "unit": "mL",
            "item": "pepper"
        },
        {
            "quantity": "50",
            "unit": "mL",
            "item": "mayonaise"
        },
        {
            "quantity": "50",
            "unit": "mL",
            "item": "mustard"
        },
        {
            "quantity": "150",
            "unit": "g",
            "item": "ham"
        },
        {
            "quantity": "150",
            "unit": "g",
            "item": "cheese"
        },
    ],
    "instr":  "spread ketchup on bread\nadd cheese and ham\ncombine slices and serve"
}


def test_modify():

  d, w = login()
  create_flow(d, w, before_img)

  d.refresh()

  # Name
  name_ele = d.find_element(By.CSS_SELECTOR, ".recipe-single-name input")
  clear(d, name_ele)
  name_ele.send_keys(after_img["name"])
  name_ele.send_keys(Keys.TAB)

  # Meal
  meal_ele = d.find_element(By.CSS_SELECTOR, ".recipe-single-meal select")
  meal_ele.send_keys(after_img["meal"])
  meal_ele.send_keys(Keys.TAB)

  # Description
  desc_ele = d.find_element(By.CSS_SELECTOR, ".recipe-single-desc textarea")
  clear(d, desc_ele)
  desc_ele.send_keys(after_img["desc"])
  desc_ele.send_keys(Keys.TAB)

  # Delete ingredient
  d.find_elements(
      By.CSS_SELECTOR, ".recipe-ingredient-delete img")[1].click()

  for i, ingre in enumerate(after_img["ingre"]):
    ingres = d.find_elements(By.CLASS_NAME, "recipe-ingredient-root")
    if i < len(ingres):
      ingre_ele = ingres[i]
    else:
      w.until(lambda x: x.find_element(
          By.CSS_SELECTOR, ".recipe-single-ingredient-add img")).click()
      w.until(lambda x: len(x.find_elements(
          By.CLASS_NAME, "recipe-ingredient-root")) > i)
      ingre_ele = d.find_elements(By.CLASS_NAME, "recipe-ingredient-root")[-1]

    quantity = ingre_ele.find_element(
        By.CSS_SELECTOR, ".recipe-ingredient-quantity input")
    unit = ingre_ele.find_element(By.CSS_SELECTOR, ".recipe-ingredient-unit")
    item = ingre_ele.find_element(
        By.CSS_SELECTOR, ".recipe-ingredient-item input")

    unit.send_keys(ingre["unit"])
    clear(d, quantity)
    quantity.send_keys(ingre["quantity"])
    clear(d, item)
    item.send_keys(ingre["item"])
    item.send_keys(Keys.ENTER)

  d.find_element(By.CSS_SELECTOR, "ol.recipe-instruction-list").click()
  instr = d.find_element(By.CSS_SELECTOR, ".recipe-instruction-list textarea")
  instr.send_keys("\ncombine slices and serve")
  instr.send_keys(Keys.TAB)

  verify(d, after_img)
