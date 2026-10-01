// Ported from the owner's Binarytree1/DiameterTree.java (isIdentical), unchanged.
class Solution {
    public boolean sameTree(TreeNode node, TreeNode subRoot) {
        if (node == null && subRoot == null) {
            return true;
        } else if (node == null || subRoot == null || node.val != subRoot.val) {
            return false;
        }
        if (!sameTree(node.left, subRoot.left)) {
            return false;
        }
        if (!sameTree(node.right, subRoot.right)) {
            return false;
        }
        return true;
    }
}
