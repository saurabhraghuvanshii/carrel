class Solution {
    public void mergeInto(int[] first, int m, int[] second) {
        int i = m - 1;
        int j = second.length - 1;
        for (int w = m + second.length - 1; j >= 0; w--) {
            if (i >= 0 && first[i] > second[j]) {
                first[w] = first[i--];
            } else {
                first[w] = second[j--];
            }
        }
    }
}
