// Ported from the owner's Binarytree1/TopView.java (lca2), unchanged; returns the node's value.
class Solution {
    public int commonAncestor(TreeNode root, int n1, int n2) {
        return lca2(root, n1, n2).val;
    }

    private TreeNode lca2(TreeNode root, int n1, int n2) {
        if (root == null || root.val == n1 || root.val == n2) {
            return root;
        }
        TreeNode leftLca = lca2(root.left, n1, n2);
        TreeNode rightLca = lca2(root.right, n1, n2);
        if (rightLca == null) {
            return leftLca;
        }
        if (leftLca == null) {
            return rightLca;
        }
        return root;
    }
}
