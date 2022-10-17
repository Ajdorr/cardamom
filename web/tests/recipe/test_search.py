from tests import login, clear
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys


def test_search():
  "Requires test_data to exist before this will run successfully"
  d, w = login()

  # Go to recipe list screen
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".workspace-menu-bar-show-btn img")).click()
  d.find_element(By.CSS_SELECTOR, "#workspace-menu-link-recipe img").click()
  w.until(lambda x: x.find_element(By.ID, "recipe-index-search-btn")).click()

  # Wait until screen loads
  name_ele = w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".recipe-search-menu-name input"))
  name_ele.send_keys("pan")
  name_ele.send_keys(Keys.ENTER)

  results_eles = w.until(lambda x: x.find_elements(
      By.XPATH, f"//div[contains(@class, 'recipe-search-list-element-name')]/span[text() = 'Pancakes']"))
  assert len(results_eles) > 0

  d.find_element(
      By.CSS_SELECTOR, ".recipe-search-menu-show-advanced img").click()

  meal_ele = d.find_element(
      By.CSS_SELECTOR, ".recipe-search-advanced-meal select")
  desc_ele = d.find_element(
      By.CSS_SELECTOR, ".recipe-search-advanced-description input")
  ingre_ele = d.find_element(
      By.CSS_SELECTOR, ".recipe-search-advanced-ingredient input")

  desc_ele.send_keys("irresistable")
  desc_ele.send_keys(Keys.ENTER)

  results_eles = w.until(lambda x: x.find_elements(
      By.XPATH, f"//div[contains(@class, 'recipe-search-list-empty')]"))
  assert len(results_eles) > 0

  clear(d, desc_ele)
  desc_ele.send_keys("breakfast")
  desc_ele.send_keys(Keys.ENTER)

  results_eles = w.until(lambda x: x.find_elements(
      By.XPATH, f"//div[contains(@class, 'recipe-search-list-element-name')]/span[text() = 'Pancakes']"))
  assert len(results_eles) > 0

  clear(d, name_ele)
  clear(d, desc_ele)
  meal_ele.send_keys("breakfast")

  results_eles = w.until(lambda x: x.find_elements(
      By.XPATH, f"//div[contains(@class, 'recipe-search-list-element-name')]/span[text() = 'Pancakes']"))
  assert len(results_eles) > 0

  ingre_ele.send_keys("e gg")
  ingre_ele.send_keys(Keys.ENTER)

  results_eles = w.until(lambda x: x.find_elements(
      By.XPATH, f"//div[contains(@class, 'recipe-search-list-empty')]"))
  assert len(results_eles) > 0

  clear(d, ingre_ele)
  ingre_ele.send_keys("egg")
  ingre_ele.send_keys(Keys.ENTER)

  results_eles = w.until(lambda x: x.find_elements(
      By.XPATH, f"//div[contains(@class, 'recipe-search-list-element-name')]/span[text() = 'Pancakes']"))
  assert len(results_eles) > 0
