class Solution {
    public String commonStart(String[] words) {
        String shared = words[0];
        for (String w : words) {
            while (!w.startsWith(shared)) {
                shared = shared.substring(0, shared.length() - 1);
            }
        }
        return shared;
    }
}
