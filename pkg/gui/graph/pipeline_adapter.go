package graph

func (g *GraphCanvas) Start(i int) {
	f := g.In[i]
	g.UpdateFunction(&f)
	g.Out[i] = f
	g.Wg.Done()
}
