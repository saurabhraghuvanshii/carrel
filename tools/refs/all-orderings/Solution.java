import java.util.*;

// Ported from the owner's backtarcking/FindPermu.java (findPermutation); it now
// collects the orderings instead of printing them.
class Solution {
    private final List<String> out = new ArrayList<>();

    public List<String> orderings(String word) {
        findPermutation(word, "");
        return out;
    }

    private void findPermutation(String str, String ans) {
        if (str.length() == 0) {
            out.add(ans);
            return;
        }
        for (int i = 0; i < str.length(); i++) {
            char curr = str.charAt(i);
            String newStr = str.substring(0, i) + str.substring(i + 1);
            findPermutation(newStr, ans + curr);
        }
    }
}
