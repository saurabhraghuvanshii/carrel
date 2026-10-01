class Solution {
    public TreeNode mirror(TreeNode root) {
        if (root == null) {
            return null;
        }
        TreeNode left = mirror(root.left);
        root.left = mirror(root.right);
        root.right = left;
        return root;
    }
}
