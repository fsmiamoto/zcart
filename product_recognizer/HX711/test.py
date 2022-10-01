from HX711 import HX711

sensor = HX711(5, 6)
sensor.set_reading_format("MSB", "MSB")
sensor.set_reference_unit(224)

print("will tare")
sensor.reset()
sensor.tare()
print("done taring")

try:
    while True:
        print(sensor.get_weight(5))
except:
    print("exiting")
    sensor.cleanup()
