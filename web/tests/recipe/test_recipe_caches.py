from typing import Any, Dict
from tests import login
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.remote.webdriver import WebDriver
from selenium.webdriver.support.wait import WebDriverWait


def test_add_missing():
  grocery_xpath_fmt = "//div[contains(@class, 'grocery-item-input')]/input[@value='%s']"
  inventory_xpath_fmt = "//div[contains(@class, 'inventory-item-input')]/input[@value='%s']"

  d, w = login()

  # Go to recipe list screen
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img")).click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-recipe img").click()

  recipes = w.until(lambda x: x.find_elements(
      By.CSS_SELECTOR, "a .recipe-list-element-name"
  ))
  recipes[0].click()

  w.until(lambda x: x.find_element(
      By.CLASS_NAME, "recipe-single-ingredient-more")).click()
  recipe_url = d.current_url
  missing_eles = w.until(lambda x: x.find_elements(
      By.CLASS_NAME, "recipe-ingredient-missing-element"))

  missing_grocery = missing_eles[0].find_element(
      By.CSS_SELECTOR, ".recipe-ingredient-missing-item span").text
  missing_eles[0].find_element(
      By.CSS_SELECTOR, ".recipe-ingredient-missing-add-grocery img").click()
  missing_inventory = missing_eles[1].find_element(
      By.CSS_SELECTOR, ".recipe-ingredient-missing-item span").text
  missing_eles[1].find_element(
      By.CSS_SELECTOR, ".recipe-ingredient-missing-add-inventory img").click()

  # Go to groceries
  d.find_element(By.CSS_SELECTOR, ".modal-panel-close-btn img").click()
  d.find_element(By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img").click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-grocery img").click()

  # Wait for load
  w.until(lambda x: x.find_elements(By.CLASS_NAME, "grocery-item-root"))
  assert len(d.find_elements(
      By.XPATH, grocery_xpath_fmt % missing_grocery)) == 1

  # Go to inventory
  d.find_element(By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img").click()
  d.find_element(
      By.CSS_SELECTOR, "#workspace-menu-link-inventory img").click()

  # wait for load
  w.until(lambda x: x.find_elements(By.CLASS_NAME, "inventory-item-root"))
  assert len(d.find_elements(
      By.XPATH, inventory_xpath_fmt % missing_inventory)) == 1

  # Go back to recipe
  d.get(recipe_url)
  w.until(lambda x: x.find_element(
      By.CLASS_NAME, "recipe-single-ingredient-more")).click()

  missing_eles = w.until(lambda x: x.find_elements(
      By.CLASS_NAME, "recipe-ingredient-missing-element"))
  missing_ingres = [
      e.find_element(
          By.CSS_SELECTOR, ".recipe-ingredient-missing-item span").text
      for e in missing_eles
  ]

  # Add all to grocery
  d.find_element(
      By.CSS_SELECTOR, ".recipe-ingredient-modal-add-all-grocery img").click()

  # Go back to groceries
  d.find_element(By.CSS_SELECTOR, ".modal-panel-close-btn img").click()
  d.find_element(By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img").click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-grocery img").click()

  w.until(lambda x: x.find_elements(By.CLASS_NAME, "grocery-item-root"))
  for ingre in missing_ingres:
    assert len(d.find_elements(By.XPATH, grocery_xpath_fmt % ingre)) == 1

  # Go back to recipe
  d.get(recipe_url)
  w.until(lambda x: x.find_element(
      By.CLASS_NAME, "recipe-single-ingredient-more")).click()

  missing_eles = w.until(lambda x: x.find_elements(
      By.CLASS_NAME, "recipe-ingredient-missing-element"))
  # Add all inventory
  d.find_element(
      By.CSS_SELECTOR, ".recipe-ingredient-modal-add-all-inventory img").click()

  # Go back to inventory and verify
  d.find_element(By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img").click()
  d.find_element(
      By.CSS_SELECTOR, "#workspace-menu-link-inventory img").click()
  w.until(lambda x: x.find_elements(By.CLASS_NAME, "inventory-item-root"))
  for ingre in missing_ingres:
    assert len(d.find_elements(By.XPATH, inventory_xpath_fmt % ingre)) == 1
