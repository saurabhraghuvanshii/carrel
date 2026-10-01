import java.util.*;

class Solution {
    public List<List<String>> cuts(String word) {
        List<List<String>> out = new ArrayList<>();
        walk(word, 0, new ArrayList<>(), out);
        return out;
    }

    private void walk(String w, int start, List<String> cur, List<List<String>> out) {
        if (start == w.length()) {
            out.add(new ArrayList<>(cur));
            return;
        }
        for (int end = start + 1; end <= w.length(); end++) {
            if (isMirror(w, start, end - 1)) {
                cur.add(w.substring(start, end));
                walk(w, end, cur, out);
                cur.remove(cur.size() - 1);
            }
        }
    }

    private boolean isMirror(String w, int i, int j) {
        while (i < j) {
            if (w.charAt(i++) != w.charAt(j--)) {
                return false;
            }
        }
        return true;
    }
}
