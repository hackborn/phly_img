# phly_img
Image processing for phly.

This is an addon pack for the [phly pipeline utility](https://github.com/hackborn/phly). See that repo for info on how to use the nodes.

## Nodes ##
* **Load Image** (phly/img/load).
    * input **file**. Supply file names to load.
    * output **0**. The loaded images, one for each file name input.
* **Save Image** (phly/img/save). Save images to a file.
    * cfg **file**. The name of the saved file. Allows variables `${src}` (the source file path), `${srcdir}` (the source directory), `${srcbase}` (the source filename base, minus the extension) and  `${srcext}` (the source extension).
    * input **0**. Image input.
    * output **0**. Image output. All input items are provided, even if the save failed.
* **Scale Image** (phly/img/scale). Resize images.
    * cfg **width**. The width of the final image. Allows variables `${w}` (source width) and `${h}` (source height) and arithmetic expressions (i.e. "(${w} * 0.5) + 10"
    * cfg **height**. The height of the final image. Allows variables `${w}` (source width) and `${h}` (source height) and arithmetic expressions (i.e. "(${w} * 0.5) + 10"
    * input **0**. Image input.
    * output **0**. The resized images.