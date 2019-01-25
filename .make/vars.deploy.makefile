FUNCTION_NAME              ?= $(MODULE_NAME)
FUNCTION_LIST              := fivecolorsd
FUNCTION_RUNTIME           := go111
FUNCTION_ENTRYPOINT        ?= Handler
FUNCTION_TIMEOUT           ?= 60
FUNCTION_MEMORY            ?= 128
define FUNCTION_ENVIRONMET =
VERSION=$(VERSION)\
PROJECT_ID=$(PROJECT_ID)
endef
comma  := ,
space  :=
space  +=
FUNCTION_ENVIRONMET  := $(subst $(space),$(comma),$(FUNCTION_ENVIRONMET))
