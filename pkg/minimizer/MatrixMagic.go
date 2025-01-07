package minimizer

// d = 1e-6
/*
def forward_dq(f: Callable[[np.ndarray], np.ndarray], x: np.ndarray, d: float) -> np.ndarray:
    """
    Approximates the Jacobian J_fx of the given function f at position x via the forward difference quotient.
    :param f: Function to estimate Jacobian for. Maps from an M-D numpy array to an N-D numpy array.
    :param x: Position to approximate Jacobian at.
    :param d: Steps size for the forward difference quotient.
    :return: A NxM numpy array containing the approximated Jacobian of f at position x.
    """
    res = np.zeros((x.shape[0], x.shape[0]))

    for i in range(0,x.shape[0]):
        off = np.zeros(x.shape[0])
        off[i] = d;
        res[i] = (f(x+off)-f(x))/d
    return res.transpose()
*/
