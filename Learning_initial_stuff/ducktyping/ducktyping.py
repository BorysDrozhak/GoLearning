
class duck(object):
    """docstring for duck"""
    def __init__(self, arg):
        super(duck, self).__init__()
        self.arg = arg
    def cry(self):
        print "cry"



class duck2(object):
    """docstring for duck"""
    def __init__(self, arg):
        super(duck, self).__init__()
        self.arg = arg
    def cry(self):
        print "cry" "cry"


if None:
    example = duck()
else:
    example = duck2()    # 

example.cry()

# ____________________________

