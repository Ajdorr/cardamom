from tests import login, clear
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.remote.webelement import WebElement

inventory_items = [
    "potato", "taco", "rice", "chicken", "cumin"
]


def test_required_data():
  d, w = login()

  # Go to inventory list screen
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, "#workspace-menu-link-inventory img")).click()

  add_item = d.find_element(
      By.CSS_SELECTOR, ".inventory-list-add-item-input input")

  for inv in inventory_items:
    add_item.send_keys(inv)
    add_item.send_keys(Keys.ENTER)

  for inv in inventory_items:
    eles = d.find_elements(
        By.XPATH,
        f"//div[contains(@class,'inventory-item-input')]/input[@value='{inv}']")
    assert len(eles) == 1


def test_basic_usage():
  d, w = login()

  # Go to inventory list screen
  w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, "#workspace-menu-link-inventory img")).click()

  add_item = d.find_element(
      By.CSS_SELECTOR, ".inventory-list-add-item-input input")

  add_item.send_keys("pepPEroni")
  add_item.send_keys(Keys.ENTER)
  add_item.send_keys("cauliflower")
  add_item.send_keys(Keys.ENTER)
  add_item.send_keys("caulIflowEr")
  add_item.send_keys(Keys.ENTER)
  add_item.send_keys("BaCon")
  add_item.send_keys(Keys.TAB)

  inv_item = "//div[contains(@class,'inventory-item-input')]/input[@value='%s']"
  assert len(d.find_elements(By.XPATH, inv_item % "pepperoni")) == 1
  assert len(d.find_elements(By.XPATH, inv_item % "cauliflower")) == 1
  assert len(d.find_elements(By.XPATH, inv_item % "bacon")) == 1

  item_ele = d.find_element(By.XPATH, inv_item % "bacon")
  clear(d, item_ele)
  item_ele.send_keys("pasta")
  item_ele.send_keys(Keys.TAB)
  assert len(d.find_elements(By.XPATH, inv_item % "bacon")) == 0
  assert len(d.find_elements(By.XPATH, inv_item % "pasta")) == 1

  item_ele = d.find_element(By.XPATH, inv_item % "cauliflower")
  clear(d, item_ele)
  item_ele.send_keys("strawberries")
  item_ele.send_keys(Keys.TAB)
  assert len(d.find_elements(By.XPATH, inv_item % "cauliflower")) == 0
  assert len(d.find_elements(By.XPATH, inv_item % "strawberries")) == 1

  d.find_element(
      By.XPATH, (inv_item % "pepperoni") +
      "/../../div[contains(@class,'inventory-item-unstock')]/img").click()
  d.find_element(
      By.XPATH, (inv_item % "pasta") +
      "/../../div[contains(@class,'inventory-item-unstock')]/img").click()
  d.find_element(
      By.XPATH, (inv_item % "strawberries") +
      "/../../div[contains(@class,'inventory-item-unstock')]/img").click()
  assert len(d.find_elements(By.XPATH, inv_item % "pepperoni")) == 0
  assert len(d.find_elements(By.XPATH, inv_item % "pasta")) == 0
  assert len(d.find_elements(By.XPATH, inv_item % "strawberries")) == 0

  d.refresh()
  add_item = w.until(lambda x: x.find_element(
      By.CSS_SELECTOR, ".inventory-list-add-item-input input"))
  assert len(d.find_elements(By.XPATH, inv_item % "pepperoni")) == 0
  assert len(d.find_elements(By.XPATH, inv_item % "pasta")) == 0
  assert len(d.find_elements(By.XPATH, inv_item % "strawberries")) == 0
