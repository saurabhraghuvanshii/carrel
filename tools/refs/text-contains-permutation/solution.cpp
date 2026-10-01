#include <string>
#include <vector>
using namespace std;

bool containsRearrangement(string& pattern, string& text) {
    size_t m = pattern.size();
    if (m > text.size()) return false;
    vector<int> need(26), have(26);
    for (size_t i = 0; i < m; i++) {
        need[pattern[i] - 'a']++;
        have[text[i] - 'a']++;
    }
    for (size_t i = m;; i++) {
        if (need == have) return true;
        if (i == text.size()) return false;
        have[text[i] - 'a']++;
        have[text[i - m] - 'a']--;
    }
}
