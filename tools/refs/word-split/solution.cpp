#include <string>
#include <vector>
using namespace std;

bool canSplitText(string& text, vector<string>& words) {
    vector<bool> ok(text.size() + 1, false);
    ok[0] = true;
    for (size_t i = 0; i < text.size(); i++) {
        if (!ok[i]) continue;
        for (auto& w : words)
            if (text.compare(i, w.size(), w) == 0 && i + w.size() <= text.size()) ok[i + w.size()] = true;
    }
    return ok[text.size()];
}
