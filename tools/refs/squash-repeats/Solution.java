// Ported from the owner's string/Removeconsecutive.java (removeCons).
class Solution {
    public String squash(String str) {
        StringBuilder result = new StringBuilder();
        result.append(str.charAt(0));
        for (int i = 1; i < str.length(); i++) {
            if (str.charAt(i) != str.charAt(i - 1)) {
                result.append(str.charAt(i));
            }
        }
        return result.toString();
    }
}
