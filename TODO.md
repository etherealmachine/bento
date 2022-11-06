# Clean up state and refactor drawing
State and drawing could ideally be moved to Draw and Update methods for the individual tags. This
would be especially nice if they used a common callback and event API. Ideally there's no real
difference between the built-in tags like input and textarea and a custom component.

# Repeating elements and templating

# Shared styles
CSS is the obviously solution, but a bit hard to fully implement in practice. Maybe just classes
would be good enough.

# Embed some assets by default
 Currently, an app will crash if, e.g. the `btn` attribute isn't specified. This is bad.
We already have two fonts embedded by default. It would be nice to also have default assets for 
frames, buttons, inputs, scrollbars, and textareas. It would also be good to make it clear how to
override the default assets when desired, either through attributes or possibly a config handed to
the UI at build time.

# Document allowed tags and style options
With pictures.

# Implement true scrollable viewports
The current scrollable textarea and paragraph is based on scrolling lines of text. It would be nice
to have a true viewport that can scroll. This would be a Box backed by an ebiten.Image larger than
the bounds of the Box itself. From there it's relatively easy to grab a Subimage ("viewport") of
this large area and draw it onto the screen bounded by the Box, with an offset for the scroll
position. Performance is a potential issue. Code complexity is another issue, but it should be
possible by essentially giving a Box a fake OuterWidth and OuterHeight, then calling Draw and
passing in a buffer that will be used to construct the viewport.
