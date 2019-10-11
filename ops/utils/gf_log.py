




from colored import fg, bg, attr



def log_fun(g, m):
    if g == "ERROR":
        print('%s%s%s:%s%s%s'%(bg('red'), g, attr(0), fg('red'), m, attr(0)))
    else:
        print('%s%s%s:%s%s%s'%(fg('yellow'), g, attr(0), fg('green'), m, attr(0)))