import java.util.*;

class Solution {
    public int kthSmallest(TreeNode root, int k) {
        ArrayDeque<TreeNode> stack = new ArrayDeque<>();
        TreeNode n = root;
        while (true) {
            while (n != null) {
                stack.push(n);
                n = n.left;
            }
            n = stack.pop();
            if (--k == 0) {
                return n.val;
            }
            n = n.right;
        }
    }
}
