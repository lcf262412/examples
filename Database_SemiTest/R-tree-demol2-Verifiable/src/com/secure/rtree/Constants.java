package com.secure.rtree;

import com.secure.ssvrtree.SSVRTNode;
import com.secure.svrtree.SVRTNode;
import com.secure.vsptree.VSPRTNode;
import com.secure.vtree.VRTNode;

public class Constants
{
    public static final int MAX_NUMBER_OF_ENTRIES_IN_NODE = 20;//结点中的最大条目数
    public static final int MIN_NUMBER_OF_ENTRIES_IN_NODE = 8;//结点中的最小条目数

    public static final int RTDataNode_Dimension = 2;

    /** Available RTree variants. */
    public static final int RTREE_LINEAR = 0;
    public static final int RTREE_QUADRATIC = 1;
    public static final int RTREE_EXPONENTIAL = 2;
    public static final int RSTAR = 3;

    public static final int NIL = -1;
    public static final RTNode NULL = null;

    public static final SVRTNode SNULL = null;
    public static final VRTNode VNULL = null;
    public static final VSPRTNode VSPNULL = null;
    public static final SSVRTNode SPNULL = null;
}