# Bento

[Demo](https://etherealmachine.github.io/bento/)

[![Go Reference](https://pkg.go.dev/badge/github.com/etherealmachine/bento.svg)](https://pkg.go.dev/github.com/etherealmachine/bento)
[![Build Status](https://github.com/etherealmachine/bento/workflows/Go/badge.svg)](https://github.com/etherealmachine/bento/actions?query=workflow%3AGo)

# An XML based UI builder for Ebiten

```
<col width="100%" height="100%" justify="center">
  <img path="profile.png"/>
  <Header>Lorem Ipsum</Header>
  <Text>
Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.
The standard chunk of Lorem Ipsum used since the 1500s is reproduced below for those interested. Sections 1.10.32 and 1.10.33 from "de Finibus Bonorum et Malorum" by Cicero are also reproduced in their exact original form, accompanied by English versions from the 1914 translation by H. Rackham.
  </Text>
  <Button>Clicks: {{.Clicks}}</Button>
</col>
```

![Screenshot at 2022-02-26 18-12-37](https://user-images.githubusercontent.com/460276/155865525-4de1fb69-803d-469a-bd55-c31ba5c38512.png)

![Screenshot at 2022-02-26 18-13-32](https://user-images.githubusercontent.com/460276/155865527-233abef0-ab24-4b46-a56b-bc2ea80a83cf.png)

## Events
```
type Event struct {
	X, Y  int       # For Click and Hover events
	Box   *Box
	Type  EventType # Click, Hover, Change, Update
	Value string    # For Change events
}
```

**onClick**
**onHover**
**onChange**
**onUpdate**
**onDraw**