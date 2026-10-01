import java.util.*;

class Autocomplete {
    static class Node {
        Node[] next = new Node[26];
        List<String> best = new ArrayList<>(); // up to three words below this node
    }

    private final Node root = new Node();
    private final Map<String, Integer> total = new HashMap<>();

    public void add(String word, int score) {
        total.merge(word, score, Integer::sum);
        Node cur = root;
        for (char ch : word.toCharArray()) {
            if (cur.next[ch - 'a'] == null) {
                cur.next[ch - 'a'] = new Node();
            }
            cur = cur.next[ch - 'a'];
            if (!cur.best.contains(word)) {
                cur.best.add(word);
            }
            cur.best.sort((a, b) -> total.get(a).equals(total.get(b)) ? a.compareTo(b) : Integer.compare(total.get(b), total.get(a)));
            if (cur.best.size() > 3) {
                cur.best.remove(3);
            }
        }
    }

    public List<String> suggest(String prefix) {
        Node cur = root;
        for (char ch : prefix.toCharArray()) {
            cur = cur.next[ch - 'a'];
            if (cur == null) {
                return new ArrayList<>();
            }
        }
        return cur.best;
    }
}
