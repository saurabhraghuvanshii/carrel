// Ported from the owner's "Binary search tree"/BST.java (createBST), with an int[] instead of an ArrayList.
class Solution {
    public TreeNode buildTree(int[] nums) {
        return createBST(nums, 0, nums.length - 1);
    }

    private TreeNode createBST(int[] arr, int st, int end) {
        if (st > end) {
            return null;
        }
        int mid = (st + end) / 2;
        TreeNode root = new TreeNode(arr[mid]);
        root.left = createBST(arr, st, mid - 1);
        root.right = createBST(arr, mid + 1, end);
        return root;
    }
}
