class Solution {
    private int best;

    public int bestPathSum(TreeNode root) {
        best = root.val;
        gain(root);
        return best;
    }

    private int gain(TreeNode n) {
        if (n == null) {
            return 0;
        }
        int left = Math.max(0, gain(n.left));
        int right = Math.max(0, gain(n.right));
        best = Math.max(best, n.val + left + right);
        return n.val + Math.max(left, right);
    }
}
