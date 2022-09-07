from HX711.HX711 import HX711


class WeightSensor:
    REFERENCE_UNIT = 224

    def __init__(self, HX711: HX711 = HX711(5, 6)):
        self.HX711 = HX711
        self.HX711.set_reading_format("MSB", "MSB")
        self.HX711.set_reference_unit(self.REFERENCE_UNIT)

    def tare(self):
        self.HX711.reset()
        self.HX711.tare()

    def get_reading(self, samples=1):
        return self.HX711.get_weight(samples)

    def cleanup(self):
        self.HX711.cleanup()
