#! /usr/bin/python3

import sys
from weight_sensor import WeightSensor
from cart_service import CartServiceClient, UpdateCartRequest, UpdateCartRequestAction
from logger import Logger

THRESHOLD = 30

def main():
    log = Logger()
    weight_sensor = WeightSensor()

    log.info("will tare sensor")
    weight_sensor.tare()
    log.info("done taring")


    current_reading, last_reading = 0.0, 0.0

    client = CartServiceClient("http://tokyo:3333")

    while True:
        try:
            current_reading, last_reading = weight_sensor.get_reading(), current_reading

            # FIXME: Stub action just for testing
            if(abs(current_reading-last_reading) >= THRESHOLD):
                log.info("will call cart service")
                request = UpdateCartRequest("1", 2, UpdateCartRequestAction.ADD)
                response = client.execute("2", request)
                log.info(f"got status {response.status_code}")

            log.info(f"sensor reading: {current_reading}")

        except (KeyboardInterrupt, SystemExit):
            log.info("received exit signal, cleaning up")
            weight_sensor.cleanup()
            sys.exit()

if __name__ == '__main__':
    main()
