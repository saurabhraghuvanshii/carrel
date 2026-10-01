import java.util.*;

class Solution {
    public List<List<Integer>> levels(TreeNode root) {
        List<List<Integer>> rows = new ArrayList<>();
        ArrayDeque<TreeNode> queue = new ArrayDeque<>();
        if (root != null) {
            queue.add(root);
        }
        while (!queue.isEmpty()) {
            List<Integer> row = new ArrayList<>();
            for (int i = queue.size(); i > 0; i--) {
                TreeNode n = queue.poll();
                row.add(n.val);
                if (n.left != null) {
                    queue.add(n.left);
                }
                if (n.right != null) {
                    queue.add(n.right);
                }
            }
            rows.add(row);
        }
        return rows;
    }
}
