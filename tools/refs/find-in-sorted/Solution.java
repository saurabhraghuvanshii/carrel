// Ported from the owner's "array codes/BinarySearch.java" (BinarySearchA), which
// printed the position; here it returns it, and -1 when the loop ends.
class Solution {
    public int find(int[] arr, int tar) {
        int n = arr.length;
        int start = 0;
        int end = n - 1;
        while (start <= end) {
            int mid = (start + end) / 2;
            if (arr[mid] == tar) {
                return mid;
            } else if (arr[mid] < tar) {
                start = mid + 1;
            } else {
                end = mid - 1;
            }
        }
        return -1;
    }
}
