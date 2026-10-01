#include <string>
using namespace std;

bool sameLetters(string& first, string& second) {
    if (first.size() != second.size()) return false;
    int count[26] = {0};
    for (size_t i = 0; i < first.size(); i++) {
        count[first[i] - 'a']++;
        count[second[i] - 'a']--;
    }
    for (int c : count) {
        if (c != 0) return false;
    }
    return true;
}
