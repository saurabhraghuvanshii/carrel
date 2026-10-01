import java.util.*;

// Ported from the owner's backtarcking/SubSet.java (findSubSets); it collects
// the subsets, and the empty subset is "" instead of the printed word null.
class Solution {
    private final List<String> out = new ArrayList<>();

    public List<String> subsets(String word) {
        findSubSets(word, "", 0);
        return out;
    }

    private void findSubSets(String str, String ans, int i) {
        if (i == str.length()) {
            out.add(ans);
            return;
        }
        findSubSets(str, ans + str.charAt(i), i + 1);
        findSubSets(str, ans, i + 1);
    }
}
