import java.util.*;

class Solution {
    public List<Integer> rightView(TreeNode root) {
        List<Integer> view = new ArrayList<>();
        ArrayDeque<TreeNode> queue = new ArrayDeque<>();
        if (root != null) {
            queue.add(root);
        }
        while (!queue.isEmpty()) {
            for (int i = queue.size(); i > 0; i--) {
                TreeNode n = queue.poll();
                if (i == 1) {
                    view.add(n.val);
                }
                if (n.left != null) {
                    queue.add(n.left);
                }
                if (n.right != null) {
                    queue.add(n.right);
                }
            }
        }
        return view;
    }
}
