from typing import Any, Dict
from tests import login
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.remote.webdriver import WebDriver
from selenium.webdriver.support.wait import WebDriverWait


def test_add_missing():
  d, w = login()

  # Go to recipe list screen
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img")).click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-recipe img").click()

  recipes = w.until(lambda x: x.find_elements(
      By.CSS_SELECTOR, "a .recipe-list-element-name"
  ))
  recipes[0].click()
  recipe_url = d.current_url

  # Open more ingredient modal panel
  w.until(lambda x: x.find_element(
      By.CLASS_NAME, "recipe-single-ingredient-more")).click()

  # find all ingredients
  ingre_eles = w.until(lambda x: x.find_elements(
      By.CSS_SELECTOR, ".recipe-ingredient-missing-element>.recipe-ingredient-missing-item>span"))
  ingredients = [ele.text for ele in ingre_eles]
  item1 = ingredients[0]
  item2 = ingredients[1]

  # Collect all groceries
  d.find_element(By.CSS_SELECTOR, ".recipe-ingredient-modal-add-all-grocery img").click()

  # Go to groceries
  d.find_element(By.CSS_SELECTOR, ".modal-panel-close-btn img").click()
  d.find_element(By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img").click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-grocery img").click()

  # Wait for load
  w.until(lambda x: x.find_elements(By.CLASS_NAME, "grocery-item-root"))
  # Check that the grocery is there
  for item in ingredients:
    assert len(d.find_elements(By.XPATH, f"//input[@value='{item}']")) == 1

  # Collect the grocery items
  d.find_element(
      By.XPATH,
      "//input[@value='%s']/../../div[contains(@class,'grocery-item-collect')]/img" % item1).click()
  d.find_element(
      By.XPATH,
      "//input[@value='%s']/../../div[contains(@class,'grocery-item-collect')]/img" % item2).click()

  # Go back to the recipe
  d.get(recipe_url)
  # Open more modal panel
  w.until(lambda x: x.find_element(
      By.CLASS_NAME, "recipe-single-ingredient-more")).click()

  # Check ingredients are the correct state
  for item in ingredients:
    grocery_btn = d.find_element(
        By.XPATH,
        f"//span[text()='{item}']/../../div[contains(@class, 'recipe-ingredient-modal-add-ingredient')]")
    if item == item1 or item == item2:
      assert "input-disabled" not in grocery_btn.get_dom_attribute("class")
    else:
      assert "input-disabled" in grocery_btn.get_dom_attribute("class")

  # recollect grocery item1
  d.find_element(
      By.XPATH,
      f"//span[text()='{item1}']/../../div[contains(@class,"
      "'recipe-ingredient-modal-add-ingredient')]/img").click()

  # Collect all inventory items
  d.find_element(By.CSS_SELECTOR, ".recipe-ingredient-modal-add-all-inventory img").click()

  # Go to groceries
  d.find_element(By.CSS_SELECTOR, ".modal-panel-close-btn img").click()
  d.find_element(By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img").click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-grocery img").click()

  # Wait for load
  w.until(lambda x: x.find_elements(By.CLASS_NAME, "grocery-item-root"))
  # Collect the grocery items
  assert len(d.find_elements(
      By.XPATH,
      f"//input[@value='{item1}']/../../div[contains(@class,'grocery-item-collect')]/img")) == 1

  # Go to inventory
  d.find_element(By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img").click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-inventory img").click()

  # wait for load
  w.until(lambda x: x.find_elements(By.CLASS_NAME, "inventory-item-root"))
  # Check all ingredients are in inventory
  for item in ingredients:
    assert len(d.find_elements(By.XPATH, f"//input[@value='{item}']")) == 1

  d.find_element(
      By.XPATH,
      f"//input[@value='{item1}']/../../div[contains(@class,'inventory-item-show-more')]/img").click()
  d.find_element(By.CSS_SELECTOR, ".inventory-modal-delete-btn img").click()

  d.find_element(
      By.XPATH,
      f"//input[@value='{item2}']/../../div[contains(@class,'inventory-item-show-more')]/img").click()
  d.find_element(By.CSS_SELECTOR, ".inventory-modal-delete-btn img").click()

  # Go back to recipe
  d.get(recipe_url)
  w.until(lambda x: x.find_element(
      By.CLASS_NAME, "recipe-single-ingredient-more")).click()

  ele = d.find_element(
      By.XPATH,
      f"//span[text()='{item1}']/../../"
      "div[contains(@class, 'recipe-ingredient-missing-add-inventory')]/img")
  assert ele.get_dom_attribute("src") == "/icons/inventory.svg"
  ele.click()

  ele = d.find_element(
      By.XPATH,
      f"//span[text()='{item2}']/../../"
      "div[contains(@class, 'recipe-ingredient-missing-add-inventory')]/img")
  assert ele.get_dom_attribute("src") == "/icons/inventory.svg"

  # Go to groceries
  d.find_element(By.CSS_SELECTOR, ".modal-panel-close-btn img").click()
  d.find_element(By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img").click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-inventory img").click()

  # wait for load
  w.until(lambda x: x.find_elements(By.CLASS_NAME, "inventory-item-root"))
  assert len(d.find_elements(By.XPATH, f"//input[@value='{item1}']")) == 1

  d.find_element(
      By.XPATH,
      f"//input[@value='{item1}']/../../div[contains(@class,'inventory-item-show-more')]/img").click()
  d.find_element(By.CSS_SELECTOR, ".inventory-modal-delete-btn img").click()
