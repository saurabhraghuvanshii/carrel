import java.util.*;

class Solution {
    public int maxDepth(TreeNode root) {
        int depth = 0;
        ArrayDeque<TreeNode> level = new ArrayDeque<>();
        if (root != null) {
            level.add(root);
        }
        while (!level.isEmpty()) {
            depth++;
            for (int i = level.size(); i > 0; i--) {
                TreeNode n = level.poll();
                if (n.left != null) {
                    level.add(n.left);
                }
                if (n.right != null) {
                    level.add(n.right);
                }
            }
        }
        return depth;
    }
}
