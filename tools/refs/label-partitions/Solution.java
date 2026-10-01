import java.util.*;

class Solution {
    public List<Integer> pieceSizes(String word) {
        int[] last = new int[26];
        for (int i = 0; i < word.length(); i++) {
            last[word.charAt(i) - 'a'] = i;
        }
        List<Integer> sizes = new ArrayList<>();
        int start = 0, end = 0;
        for (int i = 0; i < word.length(); i++) {
            end = Math.max(end, last[word.charAt(i) - 'a']);
            if (i == end) {
                sizes.add(end - start + 1);
                start = i + 1;
            }
        }
        return sizes;
    }
}
