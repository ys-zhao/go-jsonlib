package main

import (
	"fmt"

	"github.com/ys-zhao/jsonlib"
)

var xmlStr = `
<?xml version="1.0" encoding="UTF-8"?>
<breakfast_menu>
  <food>
    <name>Belgian Waffles</name>
    <price>$5.95</price>
    <description>Two of our famous Belgian Waffles with plenty of real maple syrup</description>
    <calories>650</calories>
  </food>
  <food>
    <name>Strawberry Belgian Waffles</name>
    <price>$7.95</price>
    <description>Light Belgian waffles covered with strawberries and whipped cream</description>
    <calories>900</calories>
  </food>
  <food>
    <name>Berry-Berry Belgian Waffles</name>
    <price>$8.95</price>
    <description>Light Belgian waffles covered with an assortment of fresh berries and whipped cream</description>
    <calories>900</calories>
  </food>
  <food>
    <name>French Toast</name>
    <price>$4.50</price>
    <description>Thick slices made from our homemade sourdough bread</description>
    <calories>600</calories>
  </food>
  <food>
    <name>Homestyle Breakfast</name>
    <price>$6.95</price>
    <description>Two eggs, bacon or sausage, toast, and our ever-popular hash browns</description>
    <calories>950</calories>
  </food>
</breakfast_menu>
`

func main() {
	jsonStr, _ := jsonlib.XML2JSON(xmlStr, jsonlib.X2JWithOmitRoot(false), jsonlib.X2JWithIndent(true, "", "  "))
	fmt.Println("main: json2xml with root node...")
	fmt.Println(jsonStr)

	jsonStr, _ = jsonlib.XML2JSON(xmlStr, jsonlib.X2JWithOmitRoot(true), jsonlib.X2JWithIndent(true, "", "  "))
	fmt.Println("main: json2xml without root node...")
	fmt.Println(jsonStr)
}
