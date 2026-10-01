// Ported from the owner's Binarytree1/DiameterTree.java (diameter2 and Info), unchanged. It counts
// the nodes on the path, which is how this problem defines its length.
class Solution {
    static class Info {
        int diam;
        int ht;

        Info(int diam, int ht) {
            this.diam = diam;
            this.ht = ht;
        }
    }

    public int longestPath(TreeNode root) {
        return diameter2(root).diam;
    }

    private Info diameter2(TreeNode root) {
        if (root == null) {
            return new Info(0, 0);
        }
        Info leftInfo = diameter2(root.left);
        Info rightInfo = diameter2(root.right);
        int diam = Math.max(Math.max(leftInfo.diam, rightInfo.diam), leftInfo.ht + rightInfo.ht + 1);
        int ht = Math.max(leftInfo.ht, rightInfo.ht) + 1;
        return new Info(diam, ht);
    }
}
