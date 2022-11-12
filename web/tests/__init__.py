from sys import platform
from pydantic import BaseSettings
from dotenv import load_dotenv
from typing import Optional, Tuple
from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.remote.webdriver import WebDriver
from selenium.webdriver.support.wait import WebDriverWait
from selenium.webdriver.remote.webelement import WebElement
from selenium.webdriver.common.by import By
from selenium.webdriver.common.action_chains import ActionChains
from selenium.webdriver.chrome.options import Options as ChromeOptions


class Settings(BaseSettings):
  target_host: str
  wait: float
  browser: str
  oauth_provider: str
  username: Optional[str]
  password: Optional[str]


load_dotenv(".pytest.env")
settings = Settings()


def get_driver() -> Tuple[WebDriver, WebDriverWait]:
  if settings.browser == "chrome":
    options = ChromeOptions()
    if platform == "linux":
      options.add_argument('--no-sandbox')
      options.add_argument('--headless')
      options.add_argument('--disable-extensions')
      options.add_argument('--disable-dev-shm-usage')
    driver = webdriver.Chrome(options=options)
  elif settings.browser == "safari":
    driver = webdriver.Safari()
  elif settings.browser == "firefox":
    driver = webdriver.Firefox()
  else:
    raise Exception(f"Unknown browser: {settings.browser}")

  return driver, WebDriverWait(driver, settings.wait)


def login() -> Tuple[WebDriver, WebDriverWait]:
  if settings.username is not None and settings.password is not None:
    return login_traditional()
  else:
    return login_oauth()


def login_oauth() -> Tuple[WebDriver, WebDriverWait]:

  import sys
  print(sys.path)
  driver, wait = get_driver()
  driver.get(settings.target_host)
  wait.until(lambda d: d.find_element(
      By.CLASS_NAME, "home-login-link")).click()
  wait.until(lambda d: d.find_element(
      By.ID, "auth-login-button-github")).click()

  return driver, wait


def login_traditional() -> Tuple[WebDriver, WebDriverWait]:

  driver, wait = get_driver()
  driver.get(settings.target_host)

  wait.until(lambda d: d.find_element(
      By.CLASS_NAME, "home-login-link")).click()
  # Show traditional login form
  ele = wait.until(lambda d: d.find_element(
      By.ID, "auth-login-show-traditional"))
  ele.click()

  # Authentication information
  ele = wait.until(lambda d: d.find_element(
      By.CSS_SELECTOR, "#auth-login-traditional-email input")).send_keys(settings.username)

  driver.find_element(
      By.CSS_SELECTOR,
      "#auth-login-traditional-password input").send_keys(settings.password)

  driver.find_element(
      By.CSS_SELECTOR,
      "#auth-login-traditional-submit").click()

  return driver, wait


def clear(driver: WebDriver, ele: WebElement):
  action = ActionChains(driver, 50)
  if platform == "darwin":
    action.key_down(Keys.COMMAND, element=ele)
  else:
    action.key_down(Keys.CONTROL, element=ele)
  action.send_keys_to_element(ele, "a")
  if platform == "darwin":
    action.key_up(Keys.COMMAND, element=ele)
  else:
    action.key_up(Keys.CONTROL, element=ele)

  action.perform()
  ele.send_keys(Keys.BACK_SPACE)
