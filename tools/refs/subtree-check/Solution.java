// Ported from the owner's Binarytree1/DiameterTree.java (isSubtree and isIdentical), unchanged.
class Solution {
    public boolean hasSubtree(TreeNode root, TreeNode subRoot) {
        if (root == null) {
            return false;
        }
        if (root.val == subRoot.val) {
            if (isIdentical(root, subRoot)) {
                return true;
            }
        }
        return hasSubtree(root.left, subRoot) || hasSubtree(root.right, subRoot);
    }

    private boolean isIdentical(TreeNode node, TreeNode subRoot) {
        if (node == null && subRoot == null) {
            return true;
        } else if (node == null || subRoot == null || node.val != subRoot.val) {
            return false;
        }
        return isIdentical(node.left, subRoot.left) && isIdentical(node.right, subRoot.right);
    }
}
